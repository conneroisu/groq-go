library IEEE;

use IEEE.std_logic_1164.all;

entity mux2t1 is

  port (
    i_D0, i_D1, i_S : in  std_logic;
    o_O             : out std_logic
    );

end mux2t1;

architecture behaviour of mux2t1 is
begin
  
  process (i_D0, i_D1, i_S)
  begin
    if i_s = '0' then
      o_O <= i_D0;
    else
      o_O <= i_D1;
    end if;
  end process;
  
end behaviour;
